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

    const contractId = this.roundArguments.contractId;
    const noOfAssets = this.roundArguments.assets;
    const contractFunction = "InitLedger";

    for (let i = 0; i < noOfAssets; i++) {
      const employeeId = "emp" + i.toString();
      this.assetIds.push(employeeId);
      await this.sutAdapter.invokeSmartContract(
        contractId,
        contractFunction,
        { invokerIdentity: "User1" },
        [employeeId]
      );
    }
  }

  async submitTransaction() {
    const randomId = Math.floor(Math.random() * this.roundArguments.assets);
    const employeeId = this.assetIds[randomId];

    const contractId = this.roundArguments.contractId;
    const contractFunction = "GetEmployee"; // 변경: VerifyEmployee에서 GetEmployee로 수정

    const response = await this.sutAdapter.invokeSmartContract(
      contractId,
      contractFunction,
      { invokerIdentity: "User1" },
      [employeeId]
    );

    console.log(`GetEmployee response: ${response}`);
  }

  async cleanupWorkloadModule() {
    for (const employeeId of this.assetIds) {
      const contractId = this.roundArguments.contractId;
      const contractFunction = "DeleteEmployee";

      await this.sutAdapter.invokeSmartContract(
        contractId,
        contractFunction,
        { invokerIdentity: "User1" },
        [employeeId]
      );
    }
  }
}

function createWorkloadModule() {
  return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
