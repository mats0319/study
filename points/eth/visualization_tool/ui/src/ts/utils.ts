export class library {
  id: number = -1; // unique tag
  name: string = ""; // library name
  byteCode: string = ""; // library byte code, with placeholder
  dependOn: Array<number> = []; // library dependencies
  placeholder: string = ""; // if a lib depend on this lib, replace 'placeholder' with address of this lib

  txHash: string = ""; // library deploy tx hash
  address: string = ""; // library address
}

export function replacePlaceholder(arr: Array<library>, id: number): string {
  let index = -1;
  for (let i = 0; i < arr.length; i++) {
    if (arr[i].id === id) {
      index = i;
      break;
    }
  }

  if (index === -1) {
    return "ID不存在";
  }

  const dependencies = arr[index].dependOn;
  for (let i = 0; i < dependencies.length; i++) {
    for (let j = 0; j < arr.length; j++) {
      if (arr[j].id === dependencies[i]) {
        if (arr[j].address.length < 1) {
          return "存在尚未部署的依赖";
        }

        arr[index].byteCode = arr[index].byteCode.replaceAll(arr[j].placeholder, arr[j].address.slice(2));
      }
    }
  }

  return "";
}

export function isValidAddress(value: string): boolean {
  return isValidHex(value, 40);
}

// isValidTokenID require a 64-valid-digit hex string, no matter if it has '0x' prefix
export function isValidTokenID(value: string): boolean {
  return isValidHex(value, 64);
}

function isValidHex(str: string, length: number): boolean {
  if (str.slice(0, 2) !== "0x") {
    str = "0x" + str;
  }

  if (str.length !== length+2) {
    return false;
  }

  let isValid = true;
  for (let i = 2; i < str.length; i++) {
    const char = str[i];
    if (!('0' <= char && char <= '9' || 'A' <= char && char <= 'Z' || 'a' <= char && char <= 'z')) {
      isValid = false;
      break;
    }
  }

  return isValid;
}
