package com.cbpm.backend.util;

import com.alibaba.fastjson.JSON;
import com.auth0.jwt.exceptions.InvalidClaimException;
import com.auth0.jwt.exceptions.SignatureVerificationException;
import com.auth0.jwt.exceptions.TokenExpiredException;
import com.auth0.jwt.interfaces.Claim;
import com.cbpm.backend.vo.ResponseVo;
import java.io.IOException;
import java.util.Map;
import javax.servlet.Filter;
import javax.servlet.FilterChain;
import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.annotation.WebFilter;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;

/**
 * @author Polaris
 * @version 1.0
 * @description: 拦截器类
 * @date 2022/4/9 0:51
 */

@Slf4j
@WebFilter(filterName = "JwtFilter", urlPatterns = {"/work/*","/user/register"})
public class JwtFilter implements Filter {

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
    }

    @Override
    public void destroy() {
    }

    @Override
    public void doFilter(ServletRequest servletRequest, ServletResponse servletResponse,
            FilterChain filterChain) throws ServletException, IOException {

        final HttpServletRequest request = (HttpServletRequest) servletRequest;

        final HttpServletResponse response = (HttpServletResponse) servletResponse;
        response.setContentType("application/json;charset=utf-8");
        //编码
        response.setCharacterEncoding("UTF-8");
        //读取request里面的token
        final String token = request.getHeader("Authorization");
        if ("OPTIONS".equals(request.getMethod())) {
            response.setStatus(HttpServletResponse.SC_OK);
            filterChain.doFilter(request, response);
        } else {
            if (token == null) {
                response.setStatus(401);
                response.getOutputStream()
                        .write(JSON.toJSONString(
                                ResponseVo.buildFailure("unauthorized operation", 401))
                                .getBytes());
                return;
            }
            try {
                Map<String, Claim> map = JwtUtil.parseToken(token);
                //解密token得到的信息
                String userId = map.get("userId").asString();
                String orgType = map.get("orgType").asString();
                String userEmail = map.get("email").asString();
                request.setAttribute("orgType", orgType);
                request.setAttribute("email", userEmail);
                request.setAttribute("userId", userId);
                String url = request.getRequestURI();
                System.out.println(url);
                if ("/user/register".equals(url)) {
                    if (!"admin".equals(orgType)) {
                        response.setStatus(401);
                        response.getOutputStream()
                                .write(JSON.toJSONString(
                                        ResponseVo.buildFailure("unauthorized operation", 401))
                                        .getBytes());
                        return;
                    }
                }
                filterChain.doFilter(request, response);
            } catch (SignatureVerificationException e) {
                log.info("invalid token signature for token: " + token);
                response.setStatus(403);
                response.getOutputStream()
                        .write(JSON.toJSONString(
                                ResponseVo.buildFailure("invalid token signature", 403))
                                .getBytes());
            } catch (TokenExpiredException e) {
                log.info("token has expired: " + token);
                response.setStatus(403);
                response.getOutputStream()
                        .write(JSON.toJSONString(
                                ResponseVo.buildFailure("token has expired", 403))
                                .getBytes());
            } catch (InvalidClaimException e) {
                log.info("invalid token claims for token: " + token);
                response.setStatus(403);
                response.getOutputStream()
                        .write(JSON
                                .toJSONString(ResponseVo.buildFailure("invalid token claims", 403))
                                .getBytes());
            } catch (Exception e) {
                log.info("token authentication failure: " + e.getMessage());
                response.setStatus(403);
                response.getOutputStream()
                        .write(JSON
                                .toJSONString(ResponseVo
                                        .buildFailure("token authentication failure", 403))
                                .getBytes());
            }

        }
    }

}
