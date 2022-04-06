package com.cbpm.backend.controller;

import com.cbpm.backend.serviceImpl.ApiImpl;
import com.cbpm.backend.util.LogWriter;
import com.cbpm.backend.vo.ResponseVo;
import org.springframework.web.bind.annotation.*;
import javax.annotation.Resource;
import com.alibaba.fastjson.JSONObject;

@RestController
public class ApiController {

    @Resource
    ApiImpl apiImpl;

    @Resource
    LogWriter logWriter;
    /**
     * @param jsonObject
     * @return java.lang.String
     * @author cynyard
     * @date 4/1/22
     * @update by Polaris in 1/4/2022
     */
    @PostMapping("/invoke")
    public ResponseVo invokeFunc(@RequestBody JSONObject jsonObject) throws Exception {
        ResponseVo responseVo=apiImpl.invoke(jsonObject);
        logWriter.writeLog("Request: "+jsonObject.toJSONString()+"\n"+"Response: "+responseVo.toString());
        return responseVo;
    }

    /**
     * @param jsonObject
     * @return java.lang.String
     * @author cynyard
     * @date 4/1/22
     * @update by Polaris in 1/4/2022
     */
    @PostMapping("/query")
    public ResponseVo queryFunc(@RequestBody JSONObject jsonObject) throws Exception {
        ResponseVo responseVo=apiImpl.query(jsonObject);
        logWriter.writeLog("Request: "+jsonObject.toJSONString()+"\n"+"Response: "+responseVo.toString());
        return responseVo;
    }
}

