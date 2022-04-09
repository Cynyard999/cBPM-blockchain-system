package com.cbpm.backend.controller;

/**
 * @author Polaris
 * @version 1.0
 * @description: TODO
 * @date 2022/4/5 15:18
 */

import com.cbpm.backend.serviceImpl.ApiImpl;
import com.cbpm.backend.serviceImpl.UserImpl;
import com.cbpm.backend.vo.ResponseVo;
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
    UserImpl userImpl;

    @PostMapping("/register")
    ResponseVo register(@RequestBody JSONObject jsonObject) {
        ResponseVo responseVo = userImpl.register(jsonObject);
        log.info("Request: " + jsonObject.toJSONString() + " and Response: "
                + responseVo.toString());
        return responseVo;

    }

    @PostMapping("/login")
    ResponseEntity<String> login(@RequestBody JSONObject jsonObject) {
        ResponseVo responseVo = userImpl.login(jsonObject);
        log.info("Request: " + jsonObject.toJSONString() + " and Response: "
                + responseVo.toString());
        HttpHeaders responseHeaders = new HttpHeaders();
        if (responseVo.isSuccess()) {
            responseHeaders.set("Authorization", responseVo.getMessage());
            return ResponseEntity.ok().headers(responseHeaders).body("login success");
        } else {
            return ResponseEntity.status(401).body(responseVo.getMessage());
        }
    }
}
