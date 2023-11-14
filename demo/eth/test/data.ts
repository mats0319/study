import { dataContractIns, prepareContract } from "../scripts/deploy"//@ts-ignore

describe("CashboxETH", () => {
  it("Test Data Contract", async () => {
    await prepareContract()

    await TestEventIsNotFree()
    await TestCostGas()
  })
})

async function TestEventIsNotFree(): Promise<void> {
  console.log("> --- Node: test event is not free. ---")

  const tx1 = await dataContractIns.funcWithEvent1()
  const tx1Receipt = await tx1.wait()

  const tx2 = await dataContractIns.funcWithEvent2()
  const tx2Receipt = await tx2.wait()

  const tx3 = await dataContractIns.funcWithoutEvent()
  const tx3Receipt = await tx3.wait()

  console.log("gas used with    event: (less params)", tx1Receipt.gasUsed)
  console.log("gas used with    event: (more params)", tx2Receipt.gasUsed)
  console.log("gas used without event:              ", tx3Receipt.gasUsed)

  console.log()
}

async function TestCostGas(): Promise<void> {
  console.log("> --- Node: test cost gas(calculate and storage). ---")

  const tx1 = await dataContractIns.calcMul2()
  const tx1Receipt = await tx1.wait()

  const tx2 = await dataContractIns.calcMul2Double()
  const tx2Receipt = await tx2.wait()

  const tx3 = await dataContractIns.storageNothing()
  const tx3Receipt = await tx3.wait()

  const tx4 = await dataContractIns.storageVariable()
  const tx4Receipt = await tx4.wait()

  console.log("gas used calc n*2       : ", tx1Receipt.gasUsed)
  console.log("gas used calc n*2 twice : ", tx2Receipt.gasUsed)
  console.log()
  console.log("gas used storage nothing  : ", tx3Receipt.gasUsed)
  console.log("gas used storage variable : ", tx4Receipt.gasUsed)

  console.log()
}
