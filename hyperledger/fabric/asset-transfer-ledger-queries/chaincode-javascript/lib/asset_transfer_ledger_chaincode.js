/*
 SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const { Contract } = require("fabric-contract-api");

class Chaincode extends Contract {
  // CreateAsset - create a new asset, store into chaincode state
  async CreateAsset(ctx, assetID) {
    const exists = await this.AssetExists(ctx, assetID);
    if (exists) {
      throw new Error(`The asset ${assetID} already exists`);
    }

    // ==== Create asset object and marshal to JSON ====
    let asset = {
      docType: "idCard",
      assetID: assetID,
    };

    // === Save asset to state ===
    await ctx.stub.putState(assetID, Buffer.from(JSON.stringify(asset)));
    let indexName = "cn";
    let colorNameIndexKey = await ctx.stub.createCompositeKey(indexName, [
      // asset.color,
      asset.assetID,
    ]);

    await ctx.stub.putState(colorNameIndexKey, Buffer.from("\u0000"));
  }

  // ReadAsset returns the asset stored in the world state with given id.
  async ReadAsset(ctx, id) {
    const assetJSON = await ctx.stub.getState(id); // get the asset from chaincode state
    if (!assetJSON || assetJSON.length === 0) {
      throw new Error(`Asset ${id} does not exist`);
    }

    return assetJSON.toString();
  }

  // delete - remove a asset key/value pair from state
  async DeleteAsset(ctx, id) {
    if (!id) {
      throw new Error("Asset name must not be empty");
    }

    let exists = await this.AssetExists(ctx, id);
    if (!exists) {
      throw new Error(`Asset ${id} does not exist`);
    }

    // to maintain the cn index, we need to read the asset first and get its color
    let valAsbytes = await ctx.stub.getState(id); // get the asset from chaincode state
    let jsonResp = {};
    if (!valAsbytes) {
      jsonResp.error = `Asset does not exist: ${id}`;
      throw new Error(jsonResp);
    }
    let assetJSON;
    try {
      assetJSON = JSON.parse(valAsbytes.toString());
    } catch (err) {
      jsonResp = {};
      jsonResp.error = `Failed to decode JSON of: ${id}`;
      throw new Error(jsonResp);
    }
    await ctx.stub.deleteState(id); //remove the asset from chaincode state

    // delete the index
    let indexName = "cn";
    let colorNameIndexKey = ctx.stub.createCompositeKey(indexName, [
      assetJSON.color,
      assetJSON.assetID,
    ]);
    if (!colorNameIndexKey) {
      throw new Error(" Failed to create the createCompositeKey");
    }
    //  Delete index entry to state.
    await ctx.stub.deleteState(colorNameIndexKey);
  }

  // TransferAsset transfers a asset by setting a new owner name on the asset
  async TransferAsset(ctx, assetName) {
    // async TransferAsset(ctx, assetName, newOwner) {
    let assetAsBytes = await ctx.stub.getState(assetName);
    if (!assetAsBytes || !assetAsBytes.toString()) {
      throw new Error(`Asset ${assetName} does not exist`);
    }
    let assetToTransfer = {};
    try {
      assetToTransfer = JSON.parse(assetAsBytes.toString()); //unmarshal
    } catch (err) {
      let jsonResp = {};
      jsonResp.error = "Failed to decode JSON of: " + assetName;
      throw new Error(jsonResp);
    }
    // assetToTransfer.owner = newOwner; //change the owner

    let assetJSONasBytes = Buffer.from(JSON.stringify(assetToTransfer));
    await ctx.stub.putState(assetName, assetJSONasBytes); //rewrite the asset
  }

  // Therefore, range queries are a safe option for performing update transactions based on query results.
  async GetAssetsByRange(ctx, startKey, endKey) {
    let resultsIterator = await ctx.stub.getStateByRange(startKey, endKey);
    let results = await this.GetAllResults(resultsIterator, false);

    return JSON.stringify(results);
  }

  // Example: GetStateByPartialCompositeKey/RangeQuery
  async TransferAssetByColor(ctx, color, newOwner) {
    // Query the cn index by color
    // This will execute a key range query on all keys starting with 'color'
    let coloredAssetResultsIterator =
      await ctx.stub.getStateByPartialCompositeKey("cn", [color]);

    // Iterate through result set and for each asset found, transfer to newOwner
    let responseRange = await coloredAssetResultsIterator.next();
    while (!responseRange.done) {
      if (!responseRange || !responseRange.value || !responseRange.value.key) {
        return;
      }

      let objectType;
      let attributes;
      ({ objectType, attributes } = await ctx.stub.splitCompositeKey(
        responseRange.value.key
      ));

      console.log(objectType);
      let returnedAssetName = attributes[1];

      // Now call the transfer function for the found asset.
      // Re-use the same function that is used to transfer individual assets
      await this.TransferAsset(ctx, returnedAssetName, newOwner);
      responseRange = await coloredAssetResultsIterator.next();
    }
  }

  // Example: Parameterized rich query
  async QueryAssetsByOwner(ctx, owner) {
    let queryString = {};
    queryString.selector = {};
    queryString.selector.docType = "idCard";
    // queryString.selector.owner = owner;
    return await this.GetQueryResultForQueryString(
      ctx,
      JSON.stringify(queryString)
    ); //shim.success(queryResults);
  }

  // Example: Ad hoc rich query
  // QueryAssets uses a query string to perform a query for assets.
  // Query string matching state database syntax is passed in and executed as is.
  // Supports ad hoc queries that can be defined at runtime by the client.
  // If this is not desired, follow the QueryAssetsForOwner example for parameterized queries.
  // Only available on state databases that support rich query (e.g. CouchDB)
  async QueryAssets(ctx, queryString) {
    return await this.GetQueryResultForQueryString(ctx, queryString);
  }

  // GetQueryResultForQueryString executes the passed in query string.
  // Result set is built and returned as a byte array containing the JSON results.
  async GetQueryResultForQueryString(ctx, queryString) {
    let resultsIterator = await ctx.stub.getQueryResult(queryString);
    let results = await this.GetAllResults(resultsIterator, false);

    return JSON.stringify(results);
  }

  // Example: Pagination with Range Query
  // GetAssetsByRangeWithPagination performs a range query based on the start & end key,
  // page size and a bookmark.
  // The number of fetched records will be equal to or lesser than the page size.
  // Paginated range queries are only valid for read only transactions.
  async GetAssetsByRangeWithPagination(
    ctx,
    startKey,
    endKey,
    pageSize,
    bookmark
  ) {
    const { iterator, metadata } = await ctx.stub.getStateByRangeWithPagination(
      startKey,
      endKey,
      pageSize,
      bookmark
    );
    const results = await this.GetAllResults(iterator, false);

    results.ResponseMetadata = {
      RecordsCount: metadata.fetched_records_count,
      Bookmark: metadata.bookmark,
    };
    return JSON.stringify(results);
  }

  // Example: Pagination with Ad hoc Rich Query
  // QueryAssetsWithPagination uses a query string, page size and a bookmark to perform a query
  // for assets. Query string matching state database syntax is passed in and executed as is.
  // The number of fetched records would be equal to or lesser than the specified page size.
  // Supports ad hoc queries that can be defined at runtime by the client.
  // If this is not desired, follow the QueryAssetsForOwner example for parameterized queries.
  // Only available on state databases that support rich query (e.g. CouchDB)
  // Paginated queries are only valid for read only transactions.
  async QueryAssetsWithPagination(ctx, queryString, pageSize, bookmark) {
    const { iterator, metadata } = await ctx.stub.getQueryResultWithPagination(
      queryString,
      pageSize,
      bookmark
    );
    const results = await this.GetAllResults(iterator, false);

    results.ResponseMetadata = {
      RecordsCount: metadata.fetched_records_count,
      Bookmark: metadata.bookmark,
    };

    return JSON.stringify(results);
  }

  // GetAssetHistory returns the chain of custody for an asset since issuance.
  async GetAssetHistory(ctx, assetName) {
    let resultsIterator = await ctx.stub.getHistoryForKey(assetName);
    let results = await this.GetAllResults(resultsIterator, true);

    return JSON.stringify(results);
  }

  // AssetExists returns true when asset with given ID exists in world state
  async AssetExists(ctx, assetName) {
    // ==== Check if asset already exists ====
    let assetState = await ctx.stub.getState(assetName);
    return assetState && assetState.length > 0;
  }

  async GetAllResults(iterator, isHistory) {
    let allResults = [];
    let res = await iterator.next();
    while (!res.done) {
      if (res.value && res.value.value.toString()) {
        let jsonRes = {};
        console.log(res.value.value.toString("utf8"));
        if (isHistory && isHistory === true) {
          jsonRes.TxId = res.value.tx_id;
          jsonRes.Timestamp = res.value.timestamp;
          try {
            jsonRes.Value = JSON.parse(res.value.value.toString("utf8"));
          } catch (err) {
            console.log(err);
            jsonRes.Value = res.value.value.toString("utf8");
          }
        } else {
          jsonRes.Key = res.value.key;
          try {
            jsonRes.Record = JSON.parse(res.value.value.toString("utf8"));
          } catch (err) {
            console.log(err);
            jsonRes.Record = res.value.value.toString("utf8");
          }
        }
        allResults.push(jsonRes);
      }
      res = await iterator.next();
    }
    iterator.close();
    return allResults;
  }

  // InitLedger creates sample assets in the ledger
  async InitLedger(ctx) {
    const assets = [{ assetID: "asset1" }, { assetID: "asset2" }];

    for (const asset of assets) {
      await this.CreateAsset(
        ctx,
        asset.assetID
        // asset.color,
      );
    }
  }
}

module.exports = Chaincode;
