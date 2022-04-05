package com.cbpm.backend.service;

import com.alibaba.fastjson.JSONObject;
import com.cbpm.backend.vo.ResponseVo;
/**
* @description: userInterface
 * @param
* @return:
* @author: Polaris
* @date: 2022/4/5
*/
public interface UserService {
    ResponseVo register(JSONObject jsonObject);
}
