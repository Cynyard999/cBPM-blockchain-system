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
import org.springframework.web.bind.annotation.*;
import javax.annotation.Resource;
import com.alibaba.fastjson.JSONObject;

@RestController
@RequestMapping("/user")
public class UserController {
    @Resource
    UserImpl userImpl;


    @PostMapping("/register")
    ResponseVo register(@RequestBody JSONObject jsonObject){
        return userImpl.register(jsonObject);
    }
}
