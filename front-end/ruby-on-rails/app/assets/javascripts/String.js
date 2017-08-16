
String.prototype.replaceAll = String.prototype.replaceAll || function(s, r) {
  return this.replace(new RegExp(s, 'g'), r);
};
