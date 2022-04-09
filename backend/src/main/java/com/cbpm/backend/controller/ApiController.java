package com.cbpm.backend.controller;

import com.cbpm.backend.serviceImpl.ApiImpl;
import com.cbpm.backend.util.JsonReader;
import com.cbpm.backend.vo.ResponseVo;
import javax.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;
import javax.annotation.Resource;
import com.alibaba.fastjson.JSONObject;

@Slf4j
@RestController
@RequestMapping("/work")
public class ApiController {

    @Resource
    ApiImpl apiImpl;


    /**
     * @param request
     * @return java.lang.String
     * @author cynyard
     * @date 4/1/22
     * @update by Polaris in 1/4/2022
     */
    @PostMapping("/invoke")
    public ResponseVo invokeFunc(HttpServletRequest request) throws Exception {
        String orgType = request.getAttribute("orgType").toString();
        JSONObject jsonObject = JsonReader.receivePost(request);
        jsonObject.put("orgType", orgType);
        ResponseVo responseVo = apiImpl.invoke(jsonObject);
        log.info("Request: " + jsonObject.toJSONString() + "and Response: " + responseVo
                .toString());
        return responseVo;
    }

    /**
     * @param request
     * @return java.lang.String
     * @author cynyard
     * @date 4/1/22
     * @update by Polaris in 1/4/2022
     */
    @PostMapping("/query")
    public ResponseVo queryFunc(HttpServletRequest request) throws Exception {
        String orgType = request.getAttribute("orgType").toString();
        JSONObject jsonObject = JsonReader.receivePost(request);
        jsonObject.put("orgType", orgType);
        ResponseVo responseVo = apiImpl.query(jsonObject);
        log.info("Request: " + jsonObject.toJSONString() + "and Response: " + responseVo.toString());
        return responseVo;
    }
}

