"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");
const helper = require("./helper");

class VerifyEmployeeWorkload extends WorkloadModuleBase {
  constructor() {
    super();
    this.txIndex = 0;
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
    await helper.createGateway(sutContext);
  }

  async submitTransaction() {
    // generate a unique transaction id
    this.txIndex++;
    let employeeID = "employee_" + this.workerIndex + "_" + this.txIndex;

    // Create the employee
    await helper.contract.submitTransaction("CreateEmployee", {
      DocType: "employee",
      ID: employeeID,
      DID: "",
    });

    // Then verify the employee
    await helper.contract.submitTransaction("VerifyEmployee", employeeID);
  }

  async cleanupWorkloadModule() {
    await helper.disconnectGateway();
  }
}

function createWorkloadModule() {
  return new VerifyEmployeeWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
