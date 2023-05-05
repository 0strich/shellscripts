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
      const assetID = `${this.workerIndex}_${i}`;
      console.log(`Worker ${this.workerIndex}: Creating asset ${assetID}`);
      const request = {
        contractId: this.roundArguments.contractId,
        contractFunction: "CreateAsset",
        invokerIdentity: "User1",
        contractArguments: [assetID],
        readOnly: false,
      };
      await this.sutAdapter.sendRequests(request);
      this.assetIds.push(assetID);
    }
  }

  async submitTransaction() {
    const randomId = Math.floor(Math.random() * this.roundArguments.assets);
    const assetID = this.assetIds[randomId];

    const request = {
      contractId: this.roundArguments.contractId,
      invokerIdentity: "User1",
      readOnly: false,
    };

    console.log(`Worker ${this.workerIndex}: Creating asset ${assetID}`);
    request.contractFunction = "CreateAsset";
    request.contractArguments = [assetID];

    await this.sutAdapter.sendRequests(request);
  }

  async cleanupWorkloadModule() {
    for (const assetID of this.assetIds) {
      console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
      const request = {
        contractId: this.roundArguments.contractId,
        contractFunction: "DeleteAsset",
        invokerIdentity: "User1",
        contractArguments: [assetID],
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
