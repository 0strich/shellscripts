"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");

class MyWorkload extends WorkloadModuleBase {
  constructor() {
    super();
    this.assetIds = [];
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

    for (let i = 0; i < this.roundArguments.assets; i++) {
      const employeeID = `${this.workerIndex}_${i}`;
      console.log("employeeID: ", employeeID);
      console.log(`Worker ${this.workerIndex}: Creating asset ${employeeID}`);
      const request = {
        contractId: this.roundArguments.contractId,
        contractFunction: "CreateEmployee",
        invokerIdentity: "User1",
        contractArguments: ["employee", employeeID, ""],
        readOnly: false,
      };
      await this.sutAdapter.sendRequests(request);
      this.assetIds.push(employeeID);
    }
  }

  async submitTransaction() {
    const randomId = Math.floor(Math.random() * this.roundArguments.assets);
    const employeeID = this.assetIds[randomId];

    const myArgs = {
      contractId: this.roundArguments.contractId,
      contractFunction: "GetEmployee",
      invokerIdentity: "User1",
      contractArguments: [employeeID],
      readOnly: true,
    };

    await this.sutAdapter.sendRequests(myArgs);
  }

  async cleanupWorkloadModule() {
    for (const employeeID of this.assetIds) {
      console.log(`Worker ${this.workerIndex}: Deleting asset ${employeeID}`);
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
