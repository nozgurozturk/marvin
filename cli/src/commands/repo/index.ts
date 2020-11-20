import createRepo from './create';
import deleteRepo from './delete';
import listRepo from './find';
import updateRepo from './update';

export const repository = () => {
  createRepo();
  updateRepo();
  listRepo();
  deleteRepo();
};
