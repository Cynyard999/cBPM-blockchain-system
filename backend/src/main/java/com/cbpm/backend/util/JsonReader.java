package com.cbpm.backend.util;



import com.alibaba.fastjson.JSONObject;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.UnsupportedEncodingException;
import java.nio.charset.StandardCharsets;
import javax.servlet.http.HttpServletRequest;

/**
 * @author Polaris
 * @version 1.0
 * @description: 读取request里面的json
 * @date 2022/4/9 13:19
 */
public class JsonReader {

    public static JSONObject receivePostBody(HttpServletRequest request) throws IOException, UnsupportedEncodingException {

        BufferedReader streamReader = new BufferedReader(new InputStreamReader(request.getInputStream(),
                StandardCharsets.UTF_8));
        StringBuilder responseStrBuilder = new StringBuilder();
        String inputStr;
        while ((inputStr = streamReader.readLine()) != null) {
            responseStrBuilder.append(inputStr);
        }
        return JSONObject.parseObject(responseStrBuilder.toString());
    }

}
