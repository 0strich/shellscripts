"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");

class MyWorkload extends WorkloadModuleBase {
  constructor() {
    super();
    this.employeeIDs = [];
  }

  async initializeWorkloadModule(
    workerIndex,
    totalWorkers,
    roundIndex,
    roundArguments,
    sutAdapter,
    sutContext
  ) {
    await super.initializeWorkloadModule(
      workerIndex,
      totalWorkers,
      roundIndex,
      roundArguments,
      sutAdapter,
      sutContext
    );

    for (let i = 0; i < this.roundArguments.employees; i++) {
      const employeeID = `${this.workerIndex}_${i}`;
      console.log("employeeID: ", employeeID);
      console.log(`Worker ${this.workerIndex}: Creating emp ${employeeID}`);
      const request = {
        contractId: this.roundArguments.contractId,
        contractFunction: "CreateEmployee",
        invokerIdentity: "User1",
        contractArguments: ["employee", employeeID, ""],
        readOnly: false,
      };
      await this.sutAdapter.sendRequests(request);
      this.employeeIDs.push(employeeID);
    }
  }

  async submitTransaction() {
    const randomId = Math.floor(Math.random() * this.roundArguments.employees);
    const employeeID = this.employeeIDs[randomId];

    const myArgs = {
      contractId: this.roundArguments.contractId,
      contractFunction: "GetDIDDocument",
      invokerIdentity: "User1",
      contractArguments: [employeeID],
      readOnly: true,
    };

    await this.sutAdapter.sendRequests(myArgs);
  }

  async cleanupWorkloadModule() {
    for (const employeeID of this.employeeIDs) {
      console.log(`Worker ${this.workerIndex}: Deleting emp ${employeeID}`);
      const request = {
        contractId: this.roundArguments.contractId,
        contractFunction: "DeleteEmployee",
        invokerIdentity: "User1",
        contractArguments: [employeeID],
        readOnly: false,
      };

      await this.sutAdapter.sendRequests(request);
    }
  }
}

function createWorkloadModule() {
  return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
