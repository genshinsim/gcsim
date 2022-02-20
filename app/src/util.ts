export function maxLvlToAsc(lvl: number): number {
  switch (lvl) {
    case 90:
      return 6;
    case 80:
      return 5;
    case 70:
      return 4;
    case 60:
      return 3;
    case 50:
      return 2;
    case 40:
      return 1;
    default:
      return 0;
  }
}

export function ascToMaxLvl(asc: number): number {
  switch (asc) {
    case 6:
      return 90;
    case 5:
      return 80;
    case 4:
      return 70;
    case 3:
      return 60;
    case 2:
      return 50;
    case 1:
      return 40;
    default:
      return 20;
  }
}

export function ascLvlMin(asc: number): number {
  switch (asc) {
    case 1:
      return 20;
    case 2:
      return 40;
    case 3:
      return 50;
    case 4:
      return 60;
    case 5:
      return 70;
    case 6:
      return 80;
  }
  return 1;
}

export function ascLvlMax(asc: number): number {
  switch (asc) {
    case 0:
      return 20;
    case 1:
      return 40;
    case 2:
      return 50;
    case 3:
      return 60;
    case 4:
      return 70;
    case 5:
      return 80;
    case 6:
      return 90;
  }
  return 0;
}
