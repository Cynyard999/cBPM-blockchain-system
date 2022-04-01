package com.cbpm.backend.controller;

import com.cbpm.backend.serviceImpl.ApiImpl;
import org.springframework.web.bind.annotation.*;
import javax.annotation.Resource;
import com.alibaba.fastjson.JSONObject;

@RestController
public class ApiController {

    @Resource
    ApiImpl apiImpl;


    /**
     * @author cynyard
     * @date 4/1/22
     * @param jsonObject
     * @update by Polaris in 1/4/2022
     * @return java.lang.String
     */
    @PostMapping("/invoke")
    public String invokeFunc(@RequestBody JSONObject jsonObject) throws Exception {
        return apiImpl.invoke(jsonObject);
    }

    /**
     * @author cynyard
     * @date 4/1/22
     * @param jsonObject
     * @update by Polaris in 1/4/2022
     * @return java.lang.String
     */
    @PostMapping("/query")
    public String queryFunc(@RequestBody JSONObject jsonObject) throws Exception {
        return apiImpl.query(jsonObject);
    }
}

