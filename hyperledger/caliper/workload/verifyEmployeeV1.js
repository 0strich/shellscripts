"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");
const crypto = require("crypto");

class VerifyEmployeeWorkload extends WorkloadModuleBase {
  constructor() {
    super();
    this.txIndex = 0;
    this.limitIndex = 0;
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
    this.roundArguments = roundArguments;
  }

  async submitTransaction() {
    this.txIndex++;

    // Generate DID
    let id = "emp" + this.txIndex;
    let did =
      "did:ipid:" + crypto.createHash("sha256").update(id).digest("hex");

    let txFabric = {
      contractId: "mycontract",
      contractVersion: "v0",
      contractFunction: "VerifyEmployee",
      contractArguments: [did],
      timeout: this.roundArguments.timeout
        ? parseInt(this.roundArguments.timeout)
        : 120,
    };

    if (this.txIndex === this.limitIndex) {
      return this.sutAdapter.invokeSmartContract(
        this.sutContext,
        txFabric.contractId,
        txFabric.contractVersion,
        [txFabric],
        this.timeout
      );
    } else {
      return this.sutAdapter.invokeSmartContract(
        this.sutContext,
        txFabric.contractId,
        txFabric.contractVersion,
        [txFabric]
      );
    }
  }

  async cleanupWorkloadModule() {
    const sleep = (ms) => {
      return new Promise((resolve) => setTimeout(resolve, ms));
    };
    await sleep(2000);
  }
}

function createWorkloadModule() {
  return new VerifyEmployeeWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
