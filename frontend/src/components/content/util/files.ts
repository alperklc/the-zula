/**
 * Reads a file and returns an ArrayBuffer
 * @param file
 */
export async function readFileAsync(file: Blob): Promise<ArrayBuffer> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = () => {
      if (typeof reader.result === 'string') {
        throw new Error('reader.result is expected to be an ArrayBuffer')
      }
      if (reader?.result) {
        resolve(reader?.result)
      }
    }

    reader.onerror = reject

    reader.readAsArrayBuffer(file)
  })
}
