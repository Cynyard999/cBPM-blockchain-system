package com.cbpm.backend.service;

import com.alibaba.fastjson.JSONObject;
import org.hyperledger.fabric.gateway.ContractException;

public interface ApiService {

    /**
    * @description: interface for controller
     * @param jsonObject
    * @return: java.lang.String
    * @author: Polaris
    * @date: 2022/4/1
    */
    String query(JSONObject jsonObject) throws ContractException;

    /**
    * @description: interface for controller
     * @param jsonObject
    * @return: java.lang.String
    * @author: Polaris
    * @date: 2022/4/1
    */
    String invoke(JSONObject jsonObject);

}
