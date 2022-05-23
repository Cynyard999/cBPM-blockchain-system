package com.cbpm.backend.serviceImpl;

import com.alibaba.fastjson.JSONObject;
import com.cbpm.backend.dao.UserRepository;
import com.cbpm.backend.service.UserService;
import com.cbpm.backend.util.JwtUtil;
import com.cbpm.backend.vo.ResponseVo;
import com.cbpm.backend.vo.UserVo;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

/**
 * @author Polaris
 * @version 1.0
 * @description: 实现user相关逻辑
 * @date 2022/4/5 12:49
 */
@Slf4j
@Service
public class UserImpl implements UserService {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private BCryptPasswordEncoder bCryptPasswordEncoder;

    @Override
    public ResponseVo register(JSONObject jsonObject) {
        try {
            String orgType = jsonObject.getString("orgType");
            if (orgType.length() == 0) {
                return ResponseVo.buildFailure("orgType must not be null", 400);
            }
            String userName = jsonObject.getString("userName");
            if (userName.length() == 0) {
                return ResponseVo.buildFailure("userName must not be null", 400);
            }
            if (userRepository.findByName(userName) != null) {
                return ResponseVo.buildFailure(userName + " has already been registered.", 422);
            }
            String pwd = jsonObject.getString("pwd");
            if (pwd.length() == 0) {
                return ResponseVo.buildFailure("password must not be null", 400);
            }
            pwd = bCryptPasswordEncoder.encode(pwd);
            String email = jsonObject.getString("email");
            if (userRepository.findByEmail(email) != null) {
                return ResponseVo.buildFailure(email + " has already been registered.", 422);
            }
            // TODO 获取admin用户向ca注册

            UserVo userVo = new UserVo(userName, pwd, orgType, email);
            userRepository.save(userVo);
//            String token = JwtUtil.createToken(userVo);
            JSONObject userInfo = new JSONObject();
            userInfo.put("email", userVo.getEmail());
            userInfo.put("name", userVo.getName());
            userInfo.put("orgType", userVo.getOrgType());
            return ResponseVo.buildSuccess("", userInfo);
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
        if (!bCryptPasswordEncoder.matches(password, userVo.getPassword())) {
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
}
