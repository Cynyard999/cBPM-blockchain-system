package com.cbpm.backend.serviceImpl;

import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.cbpm.backend.config.GatewayConfig;
import com.cbpm.backend.service.ApiService;
import org.apache.commons.lang3.StringUtils;
import org.hyperledger.fabric.gateway.*;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.EnumSet;
import java.util.HashMap;
import java.util.concurrent.TimeoutException;

@Service
public class ApiImpl implements ApiService {

    @Resource
    GatewayConfig gatewayConfig;

    @Override
    public String query(JSONObject jsonObject) {
        //提取信息
        String orgType=jsonObject.getString("orgType");
        String channelName=jsonObject.getString("channelName");
        String contractName=jsonObject.getString("contractName");
        try {
            //获取对应组织的gateway
            Gateway gateway = gatewayConfig.gatewayHashMap.get(orgType);
            //根据channelName获取网络中channel
            Network network = gateway.getNetwork(channelName);
            //根据contractName获取channel上的contract
            Contract contract = network.getContract(contractName);
            String functionName = jsonObject.getString("function");
            if (StringUtils.isEmpty(functionName)) {
                return "no function name";
            }
            //提取参数
            JSONArray argArray = jsonObject.getJSONArray("args");
            String[] args = new String[argArray.size()];
            for (int i = 0; i < argArray.size(); i++) {
                args[i] = argArray.getString(i);
            }
            //调用contract对应function
            byte[] queryResult = contract.evaluateTransaction(functionName, args);
            System.out.println(functionName+" "+ Arrays.deepToString(args)+" "+"success");
            return new String(queryResult, StandardCharsets.UTF_8);
        }catch (ContractException e){
            String exception=e.toString();
            System.out.println(exception);
            return exception.split(":")[1];
        }catch (GatewayRuntimeException e){
            System.out.println(e.toString());
            String [] errors=e.toString().split(":");
            return errors[errors.length-1];
        }

    }

    @Override
    public String invoke(JSONObject jsonObject) {
        String orgType=jsonObject.getString("orgType");
        String channelName=jsonObject.getString("channelName");
        String contractName=jsonObject.getString("contractName");
        try {
            Gateway gateway = gatewayConfig.gatewayHashMap.get(orgType);
            Network network = gateway.getNetwork(channelName);
            Contract contract = network.getContract(contractName);
            String functionName = jsonObject.getString("function");
            if (StringUtils.isEmpty(functionName)) {
                return "no function name";
            }
            JSONArray argArray = jsonObject.getJSONArray("args");
            String[] args = new String[argArray.size()];
            for (int i = 0; i < argArray.size(); i++) {
                args[i] = argArray.getString(i);
            }
        byte[] invokeResult = contract.createTransaction(functionName).setEndorsingPeers(
                network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                .submit(args);
            System.out.println(functionName+" "+Arrays.deepToString(args)+" "+"success");
            return new String(invokeResult, StandardCharsets.UTF_8);
        }catch (ContractException e){
            String exception=e.toString();
            System.out.println(exception);
            return exception.split(":")[1];
        }catch (GatewayRuntimeException e){
            System.out.println(e.toString());
            String [] errors=e.toString().split(":");
            return errors[errors.length-1];
        } catch (InterruptedException e) {
            e.printStackTrace();
            return e.toString();
        } catch (TimeoutException e) {
            e.printStackTrace();
            return  e.toString();
        }
    }

}
