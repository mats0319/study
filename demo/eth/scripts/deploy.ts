import { ethers } from "hardhat"
import { getAddressList } from "./address"
import { Contract } from "ethers/lib.commonjs/contract/contract"

export let dataContractIns: Contract;

export async function prepareContract(dataAddress: string = ""): Promise<void> {
  await getAddressList()

  await deployData(dataAddress)
}

async function deployData(dataAddress: string = ""): Promise<void> {
  console.log("> --- Node: on deploy cashbox eth. ---");

  if (dataAddress.length > 0) { // use exist contract
    dataContractIns = await ethers.getContractAt("Data", dataAddress as string) // contract name

    try {
      await dataContractIns.getNumber() // invoke a function in contract, test contract is exist

      console.log("> Use exist contract.")
    } catch {
      // do nothing
    }
  } else { // deploy new contract
    console.log("> Deploying new contract: ")

    // data contract
    dataContractIns = await ethers.deployContract("Data")
    await dataContractIns.waitForDeployment()
    console.log("> Data contract deployed.", await dataContractIns.getAddress(), (await dataContractIns.deploymentTransaction())?.data.length)
 }

  console.log("> Data Contract Address: ", await dataContractIns.getAddress(), dataContractIns.deploymentTransaction()?.hash)
  console.log()
}
