import { ethers } from "hardhat"
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers"

export let address_1: HardhatEthersSigner
export let address_2: HardhatEthersSigner
export let address_3: HardhatEthersSigner

export async function getAddressList(): Promise<void> {
  console.log("> --- Node: show address with balance. ---")

  let addressList: HardhatEthersSigner[] = await ethers.getSigners() // number according to 'hardhat.config.ts'
  address_1 = addressList[0]
  address_2 = addressList[1]
  address_3 = addressList[2]


  console.log("> address 1 : ", address_1.address, (await ethers.provider.getBalance(address_1.address)).toString())
  console.log("> address 2 : ", address_2.address, (await ethers.provider.getBalance(address_2.address)).toString())
  console.log("> address 3 : ", address_3.address, (await ethers.provider.getBalance(address_3.address)).toString())
  console.log()
}
