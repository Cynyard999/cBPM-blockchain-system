package com.cbpm.backend.util;

import com.auth0.jwt.interfaces.Claim;
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
@WebFilter(filterName = "JwtFilter", urlPatterns = "/work/*")
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

        final HttpServletRequest request=(HttpServletRequest) servletRequest;

        final HttpServletResponse response=(HttpServletResponse) servletResponse;
        //编码
        response.setCharacterEncoding("UTF-8");
        //读取token
        final String token=request.getHeader("authorization");
        if("OPTIONS".equals(request.getMethod())){
            response.setStatus(HttpServletResponse.SC_OK);
            filterChain.doFilter(request,response);
        }else{
            if(token==null){
                response.getWriter().write("no token");
                return;
            }

            Map<String, Claim> map=JwtUtil.parseToken(token);
            if(map==null){
                response.getWriter().write("invalid token");
                return;
            }
            //解密token得到的信息
            String userId=map.get("userId").asString();
            String orgType=map.get("orgType").asString();
            String userEmail=map.get("email").asString();
//            String password=map.get("password").asString();
            //放进requst的attribute里
            request.setAttribute("orgType",orgType);
            request.setAttribute("email",userEmail);
            request.setAttribute("userId",userId);
//            request.setAttribute("password",password);
            filterChain.doFilter(request,response);
        }
    }

}
