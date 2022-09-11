package com.cbpm.backend.controller;

/**
 * @author Polaris
 * @version 1.0
 * @description: TODO
 * @date 2022/4/5 15:18
 */

import com.cbpm.backend.service.UserService;
import com.cbpm.backend.util.JsonReader;
import com.cbpm.backend.vo.ResponseVo;
import javax.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import javax.annotation.Resource;
import com.alibaba.fastjson.JSONObject;

@Slf4j
@RestController
@RequestMapping("/user")
public class UserController {

    @Resource
    UserService userService;

    @PostMapping("/register")
    ResponseEntity<ResponseVo> register(HttpServletRequest request) throws Exception {
        // TODO 直接在参数中读取请求的body，eg: @RequestBody
        JSONObject jsonObject = JsonReader.receivePostBody(request);
        ResponseVo responseVo = userService.register(jsonObject);
        log.info("Request " + request.getRequestURI() + ": " + jsonObject.toJSONString());
        if (responseVo.isSuccess()) {
            return ResponseEntity.ok().body(responseVo);
        } else {
            return ResponseEntity.status(responseVo.getStatus()).body(responseVo);
        }

    }

    @PostMapping("/login")
    ResponseEntity<ResponseVo> login(HttpServletRequest request) throws Exception {
        JSONObject jsonObject = JsonReader.receivePostBody(request);
        ResponseVo responseVo = userService.login(jsonObject);
        log.info("Request " + request.getRequestURI() + ": " + jsonObject.toJSONString());
        HttpHeaders responseHeaders = new HttpHeaders();
        if (responseVo.isSuccess()) {
            responseHeaders.set("Authorization", responseVo.getMessage());
            return ResponseEntity.ok().headers(responseHeaders).body(responseVo);
        } else {
            return ResponseEntity.status(responseVo.getStatus()).body(responseVo);
        }
    }
}
