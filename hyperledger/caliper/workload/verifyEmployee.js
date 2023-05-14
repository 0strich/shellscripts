"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");

class VerifyEmployeeWorkload extends WorkloadModuleBase {
  constructor() {
    super();
  }

  // This function is called during the benchmark initialization.
  async initializeWorkloadModule(workerIndex, totalWorkers) {
    await super.initializeWorkloadModule(workerIndex, totalWorkers);
  }

  // This function is called in each round by the worker to send transactions to the SUT.
  async submitTransaction() {
    // generate unique employee ID
    const employeeID = `emp${this.workerIndex}_${this.txIndex}`;

    await this.sutAdapter.sendRequests({
      contractId: "DIDChaincode",
      contractFunction: "VerifyEmployee",
      invokerIdentity: "User1",
      contractArguments: [employeeID],
      readOnly: true,
    });
  }
}

function createWorkloadModule() {
  return new VerifyEmployeeWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
