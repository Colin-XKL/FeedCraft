// eslint-disable-next-line import/prefer-default-export
export const namingValidator = {
  validator: (value: string, cb: (err?: string) => void) => {
    if (!value) {
      cb();
      return;
    }
    const regex = /^[a-z0-9-_]+$/;
    if (!regex.test(value)) {
      cb(
        'Name can only contain lowercase letters, numbers, hyphens, and underscores'
      );
    } else {
      cb();
    }
  },
};
