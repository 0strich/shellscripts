"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");
const Logger = require("@hyperledger/caliper-core").CaliperUtils.getLogger(
  "my-workload.js"
);

class CreateEmployeeWorkload extends WorkloadModuleBase {
  constructor() {
    super();
    this.txIndex = 0;
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
    this.roundArguments = roundArguments;
  }

  async submitTransaction() {
    this.txIndex++;
    const employeeID = `CREATE_KEY_${this.workerIndex}_${this.txIndex}`;
    Logger.info(`Creating employeeID: ${employeeID}`);
    this.employeeIDs.push(employeeID);
    const request = {
      contractId: this.roundArguments.contractId,
      contractFunction: "CreateEmployee",
      invokerIdentity: "User1",
      contractArguments: [
        "employee",
        employeeID,
        "Korea",
        "19930621",
        "01024998196",
        "Seoul",
      ],
      readOnly: false,
    };
    await this.sutAdapter.sendRequests(request);
  }

  async cleanupWorkloadModule() {
    for (const employeeID of this.employeeIDs) {
      Logger.info(`Deleting employeeID: ${employeeID}`);
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
  return new CreateEmployeeWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
