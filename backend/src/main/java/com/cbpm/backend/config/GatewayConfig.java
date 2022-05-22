package com.cbpm.backend.config;


import lombok.extern.slf4j.Slf4j;
import org.hyperledger.fabric.gateway.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import java.io.IOException;
import java.io.Reader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.InvalidKeyException;
import java.security.PrivateKey;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.HashMap;


@Slf4j
@Configuration
public class GatewayConfig {

    /**
     * network路径
     */
    @Value("${backend.networkPath}")
    private String networkPath;
    /**
     * wallet文件夹路径
     */
    private final String walletDirectory = "wallet";
    /**
     * 网络配置文件路径
     */
    @Value("${backend.networkConfigPathFormat}")
    private String networkConfigPathFormat;
    /**
     * 四种组织名称
     */
    private final String[] orgs = {"carrier", "supplier", "manufacturer", "middleman"};

    /**
     * 四种组织msp
     */
    private final String[] orgMSPs = {"CarrierMSP", "SupplierMSP", "ManufacturerMSP",
            "MiddlemanMSP"};

    /**
     * 用户名
     */
    private final String[] orgAdminNames = {"carrier-admin", "supplier-admin", "manufacturer-admin",
            "middleman-admin"};
    /**
     * 用户证书路径后缀
     */
    private final String certificatePathSuffix = "/admin/msp/signcerts/cert.pem";
    /**
     * 用户私钥路径后缀
     */
    private String privateKeyPathSuffix = "/admin/msp/keystore/private_sk";

    /**
     * admin用户证书路径后缀
     */
    private final String certificatePathSuffixAdmin = "/msp/signcerts/cert.pem";
    /**
     * 用户私钥路径后缀
     */
    private String privateKeyPathSuufixAdmin = "/msp/keystore/private_sk";

    /**
     * @return java.util.HashMap<java.lang.String, org.hyperledger.fabric.gateway.Gateway>
     * @author cynyard
     * @date 5/22/22
     * @description 配置各个org对应的网关
     */
    @Bean(name = {"gatewayHashMap"})
    public HashMap<String, Gateway> connectGateway()
            throws IOException, InvalidKeyException, CertificateException {
        HashMap<String, Gateway> gatewayHashMap = new HashMap<>();
        //初始化网关wallet账户用于连接网络
        Wallet wallet = Wallets.newFileSystemWallet(Paths.get(this.walletDirectory));
        //初始化时将所有组织的admin用户认证信息都存进wallet中，以便后面直接调用
        for (int i = 0; i < orgs.length; i++) {
            //获取证书
            X509Certificate certificate = readX509Certificate(
                    Paths.get(this.networkPath + orgs[i] + this.certificatePathSuffix));
            //获取私钥
            PrivateKey privateKey = getPrivateKey(
                    Paths.get(this.networkPath + orgs[i] + this.privateKeyPathSuffix));
            //存进wallet
            wallet.put(orgAdminNames[i],
                    Identities.newX509Identity(orgMSPs[i], certificate, privateKey));
            //配置org-ca-admin的身份
            certificate = readX509Certificate(Paths.get(
                    this.networkPath + orgs[i] + "/" + orgs[i] + "-ca-admin"
                            + this.certificatePathSuffixAdmin));

            privateKey = getPrivateKey(Paths.get(
                    this.networkPath + orgs[i] + "/" + orgs[i] + "-ca-admin"
                            + this.privateKeyPathSuufixAdmin));
            //存进wallet
            wallet.put(orgs[i] + "-ca-admin",
                    Identities.newX509Identity(orgMSPs[i], certificate, privateKey));

            //根据connection.json 获取Fabric网络连接对象
            String networkConfigPath = String.format(this.networkConfigPathFormat, orgs[i]);
            Gateway.Builder builder = Gateway.createBuilder()
                    .identity(wallet, orgAdminNames[i])
                    .networkConfig(Paths.get(networkConfigPath));
            //把所有组织的连接的gateway存起来，以便后面直接调用
            gatewayHashMap.put(orgs[i], builder.connect());
        }
        return gatewayHashMap;
    }

    /**
     * @param certificatePath
     * @description: 获取证书
     * @return: java.security.cert.X509Certificate
     * @author: Polaris
     * @date: 2022/4/1
     */
    private static X509Certificate readX509Certificate(final Path certificatePath)
            throws IOException, CertificateException {
        try (Reader certificateReader = Files
                .newBufferedReader(certificatePath, StandardCharsets.UTF_8)) {
            return Identities.readX509Certificate(certificateReader);
        }
    }

    /**
     * @param privateKeyPath
     * @description: 获取私钥
     * @return: java.security.PrivateKey
     * @author: Polaris
     * @date: 2022/4/1
     */
    private static PrivateKey getPrivateKey(final Path privateKeyPath)
            throws IOException, InvalidKeyException {
        try (Reader privateKeyReader = Files
                .newBufferedReader(privateKeyPath, StandardCharsets.UTF_8)) {
            return Identities.readPrivateKey(privateKeyReader);
        }
    }
}

