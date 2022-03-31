package com.example.backend.controller;


import org.hyperledger.fabric.gateway.Contract;
import org.hyperledger.fabric.gateway.ContractException;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.sdk.Peer;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import javax.json.JsonArray;
import java.nio.charset.StandardCharsets;
import java.util.EnumSet;
import java.util.concurrent.TimeoutException;
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
@RestController
public class TestController {
    @Resource
    private Contract contract;

    @Resource
    private Network network;

    @Resource
    private Gateway gateway;

    @GetMapping("/ReadAsset")
    public String getUser(String assetId) throws ContractException {
        byte[] queryAResultBefore = contract.evaluateTransaction("ReadAsset",assetId);
        return new String(queryAResultBefore, StandardCharsets.UTF_8);
    }

    @GetMapping("/CreateAsset")
    public String createAsset(String assetId,String color,String size,String owner,String value) throws ContractException, InterruptedException, TimeoutException {
        byte[] invokeResult = contract.createTransaction("CreateAsset")
                .setEndorsingPeers(network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(assetId, color, size,owner,value);
        String txId = new String(invokeResult, StandardCharsets.UTF_8);
        return txId;
    }

    @GetMapping("/DeleteAsset")
    public String deleteAsset(String assetId) throws ContractException, InterruptedException, TimeoutException {
        byte[] invokeResult = contract.createTransaction("DeleteAsset")
                .setEndorsingPeers(network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(assetId);
        String txId = new String(invokeResult, StandardCharsets.UTF_8);
        return txId;
    }

    @GetMapping("/TransferAsset")
    public String transferAsset(String assetId,String newOwner) throws ContractException, InterruptedException, TimeoutException {
        System.out.println(assetId+" "+newOwner);
        byte[] invokeResult = contract.createTransaction("TransferAsset")
                .setEndorsingPeers(network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(assetId,newOwner);
        String txId = new String(invokeResult, StandardCharsets.UTF_8);
        return txId;
    }


}

