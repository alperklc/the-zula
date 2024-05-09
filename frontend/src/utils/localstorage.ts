function read<T>(key: string): T | null {
    const storedObject = localStorage.getItem(key)
    if (storedObject) {
      const object = JSON.parse(storedObject)
      return object
    }
  
    return null
  }
  
  function readMultiple<T>(keys: string[]): { [key: string]: T } {
    const result: { [key: string]: T } = {}
  
    keys.forEach((key: string) => {
      const storedObject = localStorage.getItem(key)
  
      if (storedObject) {
        const object: T = JSON.parse(storedObject)
        result[key] = object
      }
    })
  
    return result
  }
  
  function save<T>(key: string, object: T) {
    const stringifiedObject = JSON.stringify(object)
    localStorage.setItem(key, stringifiedObject)
  }
  
  function remove(key: string) {
    localStorage.removeItem(key)
  }
  
  function removeAll() {
    localStorage.clear()
  }
  
  export default { read, readMultiple, save, remove, removeAll }
  