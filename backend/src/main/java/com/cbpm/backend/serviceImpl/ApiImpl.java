package com.cbpm.backend.serviceImpl;

import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.alibaba.fastjson.TypeReference;
import com.cbpm.backend.config.GatewayConfig;
import com.cbpm.backend.service.ApiService;
import com.cbpm.backend.vo.ResponseVo;

import org.hyperledger.fabric.gateway.*;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.springframework.stereotype.Service;


import javax.annotation.Resource;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.TimeoutException;

@Service
public class ApiImpl implements ApiService {

    @Resource
    GatewayConfig gatewayConfig;

    @Override
    public ResponseVo query(JSONObject jsonObject) {
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
            //提取参数
            JSONArray argArray = jsonObject.getJSONArray("args");
            String[] args = new String[argArray.size()];
            for (int i = 0; i < argArray.size(); i++) {
                args[i] = argArray.getString(i);
            }
            //调用contract对应function
            byte[] queryResult = contract.evaluateTransaction(functionName, args);
            System.out.println(functionName+" "+ Arrays.deepToString(args)+" "+"success");
            return ResponseVo.buildSuccess(JSONObject.parseObject(new String(queryResult, StandardCharsets.UTF_8)));
        }catch (ContractException e){
            String exception=e.toString();
            System.out.println(exception);
            return ResponseVo.buildFailure(exception.split(":"));
        }catch (NullPointerException e){
            System.out.println(e.toString());
            return ResponseVo.buildFailure("orgType is null or invalid orgType");
        }catch (IllegalArgumentException e){
            String exception=e.toString();
            System.out.println(exception);
            String[] temps=exception.split(":");
            return ResponseVo.buildFailure(temps);
        }
        catch (GatewayRuntimeException e){
            System.out.println(e.toString());
            String [] errors=e.toString().split(":");
            return ResponseVo.buildFailure(errors[errors.length-1]);
        }

    }

    @Override
    public ResponseVo invoke(JSONObject jsonObject) {
        String orgType=jsonObject.getString("orgType");
        String channelName=jsonObject.getString("channelName");
        String contractName=jsonObject.getString("contractName");
        try {
            Gateway gateway = gatewayConfig.gatewayHashMap.get(orgType);
            Network network = gateway.getNetwork(channelName);
            Contract contract = network.getContract(contractName);
            String functionName = jsonObject.getString("function");
            String[] args;
            JSONObject transientDetail=jsonObject.getJSONObject("transient");
            //处理incoke含有transient得情况
            if(transientDetail!=null){
                Map<String,byte[]> argsMap=new HashMap<>();
                String str=transientDetail.toJSONString();
                String transientKey=str.substring(2,str.indexOf(":")-1);
                System.out.println("transient key: "+transientKey);
                str=str.substring(str.indexOf(":")+1,str.length()-1);
                System.out.println("transient value: "+str);
                argsMap.put(transientKey,str.getBytes());
                byte[] invokeResult = contract.createTransaction(functionName).setEndorsingPeers(
                                    network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
                                    .setTransient(argsMap)
                                    .submit();
                System.out.println("transient"+" "+functionName+" "+str+" "+"success");
                return ResponseVo.buildSuccess(JSONObject.parseObject(new String(invokeResult, StandardCharsets.UTF_8)));
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
            System.out.println(functionName+" "+Arrays.deepToString(args)+" "+"success");
            return ResponseVo.buildSuccess(JSONObject.parseObject(new String(invokeResult, StandardCharsets.UTF_8)));
        }catch (ContractException e){
            String exception=e.toString();
            System.out.println(exception);
            return ResponseVo.buildFailure(exception.split(":")[1]);
        }catch (IllegalArgumentException e){
            String exception=e.toString();
            String[] temps=exception.split(":");
            return ResponseVo.buildFailure(temps[temps.length-1]);
        }catch (NullPointerException e){
            System.out.println(e.toString());
            return ResponseVo.buildFailure("orgType is null or invalid orgType");
        } catch (GatewayRuntimeException e){
            System.out.println(e.toString());
            String [] errors=e.toString().split(":");
            return ResponseVo.buildFailure(errors[errors.length-1]);
        } catch (InterruptedException e) {
            e.printStackTrace();
            return ResponseVo.buildFailure(e.toString());
        } catch (TimeoutException e) {
            e.printStackTrace();
            return  ResponseVo.buildFailure(e.toString());
        }
    }

}
