package com.cbpm.backend.serviceImpl;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONException;
import com.alibaba.fastjson.JSONObject;
import com.cbpm.backend.service.ApiService;
import com.cbpm.backend.vo.ResponseVo;
import lombok.extern.slf4j.Slf4j;
import org.hyperledger.fabric.gateway.*;
import org.hyperledger.fabric.sdk.Peer;
import org.springframework.stereotype.Service;


import javax.annotation.Resource;
import java.nio.charset.StandardCharsets;
import java.util.*;

@Slf4j
@Service
public class ApiImpl implements ApiService {

//    @Resource
//    GatewayConfig gatewayConfig;

    @Resource(name = "gatewayHashMap")
    HashMap<String, Gateway> gatewayHashMap;


    private final String CHANNEL_NAME = "cbpmchannel";

    private final String CHAINCODE_NAME = "cbpmchaincode";

    @Override
    public ResponseVo query(JSONObject jsonObject) {
        //提取信息
        String orgType = jsonObject.getString("orgType");
        try {
            //获取对应组织的gateway
            Gateway gateway = gatewayHashMap.get(orgType);
            //根据channelName获取网络中channel
            Network network = gateway.getNetwork(CHANNEL_NAME);
            //根据contractName获取channel上的contract
            Contract contract = network.getContract(CHAINCODE_NAME);
            String functionName = jsonObject.getString("function");
            //提取参数
            JSONArray argArray = jsonObject.getJSONArray("args");
            String[] args = new String[argArray.size()];
            for (int i = 0; i < argArray.size(); i++) {
                args[i] = argArray.getString(i);
            }
            //调用contract对应function
            byte[] queryResult = contract.evaluateTransaction(functionName, args);
            return ResponseVo
                    .buildSuccess(null,
                            JSON.parse(new String(queryResult, StandardCharsets.UTF_8)));
        } catch (JSONException e) {
            return ResponseVo.buildFailure("Fail to parse query result: " + e.getMessage(), 500);
        } catch (ContractException e) {
            return ResponseVo
                    .buildFailure("Chaincode function querying failed: " + e.getMessage(), 500);
        } catch (IllegalArgumentException e) {
            return ResponseVo.buildFailure("Illegal argument: " + e.getMessage(), 400);
        } catch (NullPointerException e) {
            return ResponseVo.buildFailure("Null pointer detected: " + e.getMessage(), 500);
        } catch (GatewayRuntimeException e) {
            return ResponseVo.buildFailure("Runtime limit exceed: " + e.getMessage(), 504);
        } catch (Exception e) {
            return ResponseVo.buildFailure(e.getMessage(), 500);
        }
    }

    @Override
    public ResponseVo invoke(JSONObject jsonObject) {
        String orgType = jsonObject.getString("orgType");
        try {
            Gateway gateway = gatewayHashMap.get(orgType);
            Network network = gateway.getNetwork(CHANNEL_NAME);
            Contract contract = network.getContract(CHAINCODE_NAME);
            String functionName = jsonObject.getString("function");
            String[] args;
            JSONObject transientDetail = jsonObject.getJSONObject("transient");
            //处理incoke含有transient得情况
            if (transientDetail != null) {
                Map<String, byte[]> argsMap = new HashMap<>();
                String str = transientDetail.toJSONString();
                String transientKey = str.substring(2, str.indexOf(":") - 1);
                str = str.substring(str.indexOf(":") + 1, str.length() - 1);
                argsMap.put(transientKey, str.getBytes());
                byte[] invokeResult = contract.createTransaction(functionName).setEndorsingPeers(
                        network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                        .setTransient(argsMap).submit();
                return ResponseVo
                        .buildSuccess("",
                                JSON.parse(new String(invokeResult, StandardCharsets.UTF_8)));
            }
            //无transient情况
            JSONArray argArray = jsonObject.getJSONArray("args");
            args = new String[argArray.size()];
            for (int i = 0; i < argArray.size(); i++) {
                args[i] = argArray.getString(i);
            }
            byte[] invokeResult = contract.createTransaction(functionName).setEndorsingPeers(
                    network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                    .submit(args);
            return ResponseVo
                    .buildSuccess("", JSON.parse(new String(invokeResult, StandardCharsets.UTF_8)));
        } catch (JSONException e) {
            return ResponseVo.buildFailure("Fail to parse invoke result: " + e.getMessage(), 500);
        } catch (ContractException e) {
            return ResponseVo
                    .buildFailure("Chaincode function invoking failed: " + e.getMessage(), 500);
        } catch (IllegalArgumentException e) {
            return ResponseVo.buildFailure("Illegal argument: " + e.getMessage(), 400);
        } catch (NullPointerException e) {
            return ResponseVo.buildFailure("Null pointer detected: " + e.getMessage(), 500);
        } catch (GatewayRuntimeException e) {
            return ResponseVo.buildFailure("Runtime limit exceed: " + e.getMessage(), 504);
        } catch (Exception e) {
            return ResponseVo.buildFailure(e.getMessage(), 500);
        }
    }

}
