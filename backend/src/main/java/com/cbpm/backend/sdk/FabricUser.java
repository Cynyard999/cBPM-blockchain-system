package com.cbpm.backend.sdk;

import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;

import java.io.Serializable;
import java.util.Set;

/**
 * @author Polaris
 * @version 1.0
 * @description: TODO
 * @date 2022/4/4 1:35
 */
public class FabricUser implements User, Serializable {
    private static final long serialVersionUID = 1L;
    private String name;
    private Set<String> roles;
    private String account;
    private String affiliation;
    private String organization;
    private String enrollmentSecret;
    Enrollment enrollment = null;
    public FabricUser(String name, String org) {
        this.name = name;
        this.organization = org;
    }
    @Override
    public String getName() {
        return this.name;
    }

    @Override
    public Set<String> getRoles() {
        return this.roles;
    }

    public void setRoles(Set<String> roles) {
        this.roles = roles;
    }

    @Override
    public String getAccount() {
        return this.account;
    }

    public void setAccount(String account) {
        this.account = account;
    }

    @Override
    public String getAffiliation() {
        return this.affiliation;
    }
    public void setAffiliation(String af) {
        this.affiliation = af;
    }

    @Override
    public Enrollment getEnrollment() {
        return this.enrollment;
    }

    public boolean isEnrolled() {
        return this.enrollment != null;
    }

    public String getEnrollmentSecret() {
        return this.enrollmentSecret;
    }

    public void setEnrollmentSecret(String es) {
        this.enrollmentSecret = es;
    }

    public void setEnrollment(Enrollment e) {
        this.enrollment = e;
    }
    String mspId;
    @Override
    public String getMspId() {
        return this.mspId;
    }
    public void setMspId(String mspId) {
        this.mspId = mspId;
    }
}
