package com.cbpm.backend.serviceImpl;

import com.alibaba.fastjson.JSONObject;
import com.cbpm.backend.dao.UserRepository;
import com.cbpm.backend.sdk.EnrollmentImpl;
import com.cbpm.backend.sdk.FabricUser;
import com.cbpm.backend.service.UserService;
import com.cbpm.backend.util.JwtUtil;
import com.cbpm.backend.vo.ResponseVo;
import com.cbpm.backend.vo.UserVo;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.hyperledger.fabric.gateway.Wallet;
import org.hyperledger.fabric.gateway.Wallets;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric.sdk.security.CryptoSuiteFactory;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.hyperledger.fabric_ca.sdk.RegistrationRequest;
import org.hyperledger.fabric.gateway.Identities;
import org.hyperledger.fabric.gateway.Identity;
import org.hyperledger.fabric.gateway.X509Identity;

import java.io.*;
import java.nio.file.Paths;
import java.security.PrivateKey;
import java.security.Security;
import java.util.HashMap;
import java.util.Properties;

/**
 * @author Polaris
 * @version 1.0
 * @description: 实现user相关逻辑
 * @date 2022/4/5 12:49
 */
@Slf4j
@Service
public class UserImpl implements UserService {

    @Value("${backend.networkPath}")
    private String networkPath;

    @Autowired
    private UserRepository userRepository;

    public HashMap<String, String> orgMSPHash = new HashMap<String, String>() {
        {
            put("carrier", "CarrierMSP");
            put("supplier", "SupplierMSP");
            put("middleman", "MiddlemanMSP");
            put("manufacturer", "ManufacturerMSP");
        }
    };

    public HashMap<String, String> orgAdminUrlHash = new HashMap<String, String>() {
        {
            put("carrier", "https://0.0.0.0:7056");
            put("supplier", "https://0.0.0.0:7055");
            put("middleman", "https://0.0.0.0:7057");
            put("manufacturer", "https://0.0.0.0:7054");
        }
    };

    public HashMap<String, String> orgNameHash = new HashMap<String, String>() {
        {
            put("carrier", "CarrierOrg");
            put("supplier", "SupplierOrg");
            put("middleman", "MiddlemanOrg");
            put("manufacturer", "ManufacturerOrg");
        }
    };

    @Override
    public ResponseVo register(JSONObject jsonObject) {
        try {
            Wallet wallet = Wallets.newFileSystemWallet(Paths.get("wallet"));
            String orgType = jsonObject.getString("orgType");
            if (orgType.length() == 0) {
                return ResponseVo.buildFailure("orgType must not be null", 400);
            }
            String userName = jsonObject.getString("userName");
            if (userName.length() == 0) {
                return ResponseVo.buildFailure("userName must not be null", 400);
            }
            if (wallet.get(orgType + "-" + userName) != null) {
                return ResponseVo.buildFailure("user " + userName + " is already registered", 422);
            }

            String pwd = jsonObject.getString("pwd");
            if (pwd.length() == 0) {
                return ResponseVo.buildFailure("password must not be null", 400);
            }
            String email = jsonObject.getString("email");
            if (userRepository.findByEmail(email) != null) {
                return ResponseVo.buildFailure(email + " has already been registered.", 422);
            }
            UserVo userVo = new UserVo(userName, pwd, orgType, email);
            userName = orgType + "-" + userName;
            //设置安全属性
            Security.addProvider(new org.bouncycastle.jce.provider.BouncyCastleProvider());
            Properties props = new Properties();
            //设置对应org的ca证书
            props.put("pemFile", this.networkPath + orgType + "/ca-cert.pem");
            props.put("allowAllHostNames", "true");
            //创建caclient实例
            HFCAClient caClient = HFCAClient.createNewInstance(orgAdminUrlHash.get(orgType), props);
            //设置加密
            CryptoSuite cryptoSuite = CryptoSuiteFactory.getDefault().getCryptoSuite();
            //获取对应org的ca-admin用户，只有该用户有权register新user
            User userAdmin = getFabricUserLocal(orgType + "-ca-admin", orgNameHash.get(orgType),
                    orgMSPHash.get(orgType));
            caClient.setCryptoSuite(cryptoSuite);
            //设置创建新user的请求
            RegistrationRequest registrationRequest = new RegistrationRequest(userName);
            //设置所属affiliation
            registrationRequest.setAffiliation(orgType);
            //设置useraName
            registrationRequest.setEnrollmentID(userName);
            //user类型暂时无法调用链码，默认值为client
//            registrationRequest.setType("user");
            //pwd若不设置则系统会自动生成一个乱序密码
            registrationRequest.setSecret(pwd);
            registrationRequest.setMaxEnrollments(-1);
            //设置user属性，可有可无
//            registrationRequest.addAttribute(new Attribute("attr1", "value1"));	//user-defined attribute
            String enrollmentSecret = caClient.register(registrationRequest, userAdmin);
            //密码若与系统register返回的不一样则报错
            if (!pwd.equals(enrollmentSecret)) {
                return ResponseVo.buildFailure("get wrong pwd from ca-server", 500);
            } else {
//                System.out.println("get right pwd from ca-server");
            }
            //向系统enroll，获取证书和私钥
            Enrollment enrollment = caClient.enroll(userName, enrollmentSecret);
            Identity newUser = Identities.newX509Identity(orgMSPHash.get(orgType), enrollment);
            //将新user存进wallet
            wallet.put(userName, newUser);
            userRepository.save(userVo);
            String token = JwtUtil.createToken(userVo);
            JSONObject userInfo = new JSONObject();
            userInfo.put("email", userVo.getEmail());
            userInfo.put("name", userVo.getName());
            userInfo.put("orgType", userVo.getOrgType());
            return ResponseVo.buildSuccess(token, userInfo);
        } catch (Exception e) {
            log.info("register failure: " + e.getMessage());
            return ResponseVo.buildFailure(e.getMessage(), 500);
        }


    }


    @Override
    public ResponseVo login(JSONObject jsonObject) {
        String email = jsonObject.getString("email");
        if (email.length() == 0) {
            return ResponseVo.buildFailure("email must not be null.", 400);
        }

        String password = jsonObject.getString("pwd");
        if (password.length() == 0) {
            return ResponseVo.buildFailure("password must not be null.", 400);
        }
        UserVo userVo = userRepository.findByEmail(email);
        if (userVo == null) {
            return ResponseVo
                    .buildFailure("this email has not been registered, please register first.",
                            400);
        }
        if (!userVo.getPassword().equals(password)) {
            return ResponseVo.buildFailure("wrong email or password.", 400);
        } else {
            String token = JwtUtil.createToken(userVo);
            JSONObject userInfo = new JSONObject();
            userInfo.put("email", userVo.getEmail());
            userInfo.put("name", userVo.getName());
            userInfo.put("orgType", userVo.getOrgType());
            return ResponseVo.buildSuccess(token, userInfo);
        }
    }


    /**
     * @param username
     * @param org
     * @param orgId
     * @description: 获取org对应的ca-admin FabricUser对象
     * @return: org.hyperledger.fabric.sdk.User
     * @author: Polaris
     * @date: 2022/4/5
     */
    private static User getFabricUserLocal(String username, String org, String orgId)
            throws Exception {
        FabricUser user = new FabricUser(username, org);
        user.setMspId(orgId);
        Wallet wallet = Wallets.newFileSystemWallet(Paths.get("wallet"));
        X509Identity adminIdentity = (X509Identity) wallet.get(username);
        String certificate = Identities.toPemString(adminIdentity.getCertificate());
        PrivateKey pk = adminIdentity.getPrivateKey();
        EnrollmentImpl enrollement = new EnrollmentImpl(pk, certificate);
        user.setEnrollment(enrollement);
        return user;
    }
}
