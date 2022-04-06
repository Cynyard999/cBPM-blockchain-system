package com.cbpm.backend.controller;

/**
 * @author Polaris
 * @version 1.0
 * @description: TODO
 * @date 2022/4/5 15:18
 */

import com.cbpm.backend.serviceImpl.ApiImpl;
import com.cbpm.backend.serviceImpl.UserImpl;
import com.cbpm.backend.util.LogWriter;
import com.cbpm.backend.vo.ResponseVo;
import org.springframework.web.bind.annotation.*;
import javax.annotation.Resource;
import com.alibaba.fastjson.JSONObject;

@RestController
@RequestMapping("/user")
public class UserController {

    @Resource
    UserImpl userImpl;

    @Resource
    LogWriter logWriter;

    @PostMapping("/register")
    ResponseVo register(@RequestBody JSONObject jsonObject) {
        ResponseVo responseVo=userImpl.register(jsonObject);
        logWriter.writeLog("Request: "+jsonObject.toJSONString()+"\n"+"Response: "+responseVo.toString());
        return responseVo;

    }

    @PostMapping("/login")
    ResponseVo login(@RequestBody JSONObject jsonObject) {
        ResponseVo responseVo=userImpl.login(jsonObject);
        logWriter.writeLog("Request: "+jsonObject.toJSONString()+"\n"+"Response: "+responseVo.toString());
        return responseVo;
    }
}
