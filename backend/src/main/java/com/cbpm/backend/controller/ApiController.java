package com.cbpm.backend.controller;


import javax.json.JsonArray;
import org.apache.commons.lang3.StringUtils;
import org.hyperledger.fabric.gateway.Contract;
import org.hyperledger.fabric.gateway.ContractException;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.sdk.Peer;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import java.nio.charset.StandardCharsets;
import java.util.EnumSet;
import java.util.concurrent.TimeoutException;
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;

@RestController
public class ApiController {

    @Resource
    private Contract contract;

    @Resource
    private Network network;

    @Resource
    private Gateway gateway;

    /*
     * @author cynyard
     * @date 4/1/22
     * @param jsonObject
     * @return java.lang.String
     */
    @PostMapping("/invoke")
    public String invokeFunc(@RequestBody JSONObject jsonObject) throws Exception {
        String functionName = jsonObject.getString("function");
        if (StringUtils.isNoneEmpty(functionName)) {
            return "no function name";
        }
        JSONArray argArray = jsonObject.getJSONArray("args");
        String[] args = new String[argArray.size()];
        for (int i = 0; i < argArray.size(); i++) {
            args[i] = argArray.getJSONObject(i).toString();
        }
        byte[] invokeResult = contract.createTransaction(functionName).setEndorsingPeers(
                network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(args);
        return new String(invokeResult, StandardCharsets.UTF_8);
    }

    /*
     * @author cynyard
     * @date 4/1/22
     * @param jsonObject
     * @return java.lang.String
     */
    @PostMapping("/query")
    public String queryFunc(@RequestBody JSONObject jsonObject) throws Exception {
        String functionName = jsonObject.getString("function");
        if (StringUtils.isNoneEmpty(functionName)) {
            return "no function name";
        }
        JSONArray argArray = jsonObject.getJSONArray("args");
        String[] args = new String[argArray.size()];
        for (int i = 0; i < argArray.size(); i++) {
            args[i] = argArray.getJSONObject(i).toString();
        }
        byte[] queryAResultBefore = contract.evaluateTransaction(functionName, args);
        return new String(queryAResultBefore, StandardCharsets.UTF_8);
    }


    @GetMapping("/ReadAsset")
    public String readAsset(String assetId) throws ContractException {
        byte[] queryAResultBefore = contract.evaluateTransaction("ReadAsset", assetId);
        return new String(queryAResultBefore, StandardCharsets.UTF_8);
    }

    @GetMapping("/CreateAsset")
    public String createAsset(String assetId, String color, String size, String owner, String value)
            throws ContractException, InterruptedException, TimeoutException {
        byte[] invokeResult = contract.createTransaction("CreateAsset")
                .setEndorsingPeers(
                        network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(assetId, color, size, owner, value);
        String txId = new String(invokeResult, StandardCharsets.UTF_8);
        return txId;
    }

    @GetMapping("/DeleteAsset")
    public String deleteAsset(String assetId)
            throws ContractException, InterruptedException, TimeoutException {
        byte[] invokeResult = contract.createTransaction("DeleteAsset")
                .setEndorsingPeers(
                        network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(assetId);
        String txId = new String(invokeResult, StandardCharsets.UTF_8);
        return txId;
    }

    @GetMapping("/TransferAsset")
    public String transferAsset(String assetId, String newOwner)
            throws ContractException, InterruptedException, TimeoutException {
        byte[] invokeResult = contract.createTransaction("TransferAsset")
                .setEndorsingPeers(
                        network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(assetId, newOwner);
        String txId = new String(invokeResult, StandardCharsets.UTF_8);
        return txId;
    }
}

