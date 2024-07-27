// eslint-disable-next-line import/prefer-default-export
export const namingValidator = {
  validator: (value: string, cb: (err?: string) => void) => {
    const regex = /^[a-z0-9-]+$/;
    if (!regex.test(value)) {
      cb('Name can only contain lowercase letters, numbers, and hyphens');
    } else {
      cb();
    }
  },
};
