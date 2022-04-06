package com.cbpm.backend.service;

import com.alibaba.fastjson.JSONObject;
import com.cbpm.backend.vo.ResponseVo;
import org.hyperledger.fabric.gateway.ContractException;

public interface ApiService {

    /**
     * @param jsonObject
     * @description: interface for controller
     * @return: java.lang.String
     * @author: Polaris
     * @date: 2022/4/1
     */
    ResponseVo query(JSONObject jsonObject) throws ContractException;

    /**
     * @param jsonObject
     * @description: interface for controller
     * @return: java.lang.String
     * @author: Polaris
     * @date: 2022/4/1
     */
    ResponseVo invoke(JSONObject jsonObject);

}
