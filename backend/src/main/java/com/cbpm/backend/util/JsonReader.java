package com.cbpm.backend.util;



import com.alibaba.fastjson.JSONObject;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.UnsupportedEncodingException;
import javax.annotation.Resource;
import javax.servlet.http.HttpServletRequest;
import org.springframework.context.annotation.Configuration;

/**
 * @author Polaris
 * @version 1.0
 * @description: 读取request里面的json
 * @date 2022/4/9 13:19
 */
@Configuration
public class JsonReader {

    public static JSONObject receivePost(HttpServletRequest request) throws IOException, UnsupportedEncodingException {

        BufferedReader streamReader = new BufferedReader(new InputStreamReader(request.getInputStream(), "UTF-8"));
        StringBuilder responseStrBuilder = new StringBuilder();
        String inputStr;
        while ((inputStr = streamReader.readLine()) != null) {
            responseStrBuilder.append(inputStr);
        }
        JSONObject jsonObject = JSONObject.parseObject(responseStrBuilder.toString());
        return jsonObject;
    }

}
