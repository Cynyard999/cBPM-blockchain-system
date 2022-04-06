package com.cbpm.backend.sdk;

import org.hyperledger.fabric.sdk.Enrollment;

import java.io.Serializable;
import java.security.PrivateKey;

/**
 * @author Polaris
 * @version 1.0
 * @description: TODO
 * @date 2022/4/4 1:36
 */
public class EnrollmentImpl implements Enrollment, Serializable {

    private static final long serialVersionUID = 1L;
    private final PrivateKey privateKey;
    private final String certificate;

    public EnrollmentImpl(PrivateKey pk, String c) {
        this.privateKey = pk;
        this.certificate = c;
    }

    @Override
    public PrivateKey getKey() {
        return privateKey;
    }

    @Override
    public String getCert() {
        return certificate;
    }
}
