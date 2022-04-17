package com.cbpm.backend.util;

import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTVerifier;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.Claim;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.cbpm.backend.vo.UserVo;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * @author Polaris
 * @version 1.0
 * @description: Jwt工具类
 * @date 2022/4/9 0:22
 */
public class JwtUtil {

    private static final Logger logger = LoggerFactory.getLogger(JwtUtil.class);
    //盐随机写的一个
    private static final String SERECT = "asfakjGUIYFGFASKJFGhasdgjkAFHzlksh";
    //单位为秒
    private static final Long EXPIRATION = 24*60*60L;

    /**
     * @param userVo
     * @description: 生成token
     * @return: java.lang.String
     * @author: Polaris
     * @date: 2022/4/9
     */
    public static String createToken(UserVo userVo) {
        Date expireDate = new Date(System.currentTimeMillis() + EXPIRATION * 1000);
        Map<String, Object> map = new HashMap<>();
        //设置加密方式
        map.put("alg", "H256");
        //类型
        map.put("typ", "JWT");
        //将user相关信息加密生成token
        String token = JWT.create()
                .withHeader(map)
                .withClaim("userId", userVo.getId())
                .withClaim("orgType", userVo.getOrgType())
                .withClaim("email", userVo.getEmail())
//                .withClaim("password",userVo.getPassword())
                .withExpiresAt(expireDate)
                .withIssuedAt(new Date())
                .sign(Algorithm.HMAC256(SERECT));
        return token;
    }

    /**
     * @param token
     * @description: 解析token, 抛出异常
     * @return: java.util.Map<java.lang.String, com.auth0.jwt.interfaces.Claim>
     * @author: Polaris
     * @date: 2022/4/9
     * @updator: cynyard
     */
    public static Map<String, Claim> parseToken(String token) throws Exception {
        JWTVerifier jwtVerifier = JWT.require(Algorithm.HMAC256(SERECT)).build();
        DecodedJWT  decodedJWT= jwtVerifier.verify(token);
        Map<String, Claim> claimMap = decodedJWT.getClaims();
        return claimMap;
    }
}
