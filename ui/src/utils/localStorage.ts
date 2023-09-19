// This variable is used to cache the result of the localStorageAvailable function
// so that we don't have to check for localStorage availability every time we want
// to use it, reducing writes to disk. If the localStorage becomes available, a page
// refresh is required to update the value of this variable.
let cachedAvailable: boolean | null = null;

export const localStorageAvailable = () => {
  if (cachedAvailable !== null) {
    return cachedAvailable;
  }
  try {
    const storage = window["localStorage"];
    const x = "__storage_test__";
    storage.setItem(x, x);
    storage.removeItem(x);
    cachedAvailable = true;
    return true;
  } catch (e) {
    cachedAvailable = false;
    return false;
  }
};

export const getFromLocalStorage = (key: string, defaultValue: any) => {
  if (localStorageAvailable()) {
    const value = localStorage.getItem(key);
    if (value) {
      return value;
    }
  }
  return defaultValue;
};

export const saveToLocalStorage = (key: string, value: any) => {
  if (localStorageAvailable()) {
    localStorage.setItem(key, value);
  }
};

export const removeFromLocalStorage = (key: string) => {
  if (localStorageAvailable()) {
    localStorage.removeItem(key);
  }
};
