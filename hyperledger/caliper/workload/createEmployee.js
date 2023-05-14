"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");
const helper = require("./helper");

class CreateEmployeeWorkload extends WorkloadModuleBase {
  constructor() {
    super();
    this.txIndex = 0;
  }

  async initializeWorkloadModule(workerIndex, totalWorkers) {
    const admin = await helper.getOrgAdmin(
      helper.getOrganization(),
      this.sutAdapter
    );
    this.contract = admin.getContract("DIDChaincode");
  }

  async submitTransaction() {
    this.txIndex++;
    const id = `EMP${this.workerIndex}_${this.txIndex}`;
    const docType = "employee";
    await this.contract.submitTransaction("CreateEmployee", docType, id);
  }

  async cleanupWorkloadModule() {
    // No-op
  }
}

function createWorkloadModule() {
  return new CreateEmployeeWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
