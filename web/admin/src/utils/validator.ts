export const namingValidator = {
  validator: (value: string, cb: (err?: string) => void) => {
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
