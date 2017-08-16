import { PharmacoDBPage } from './app.po';

describe('pharmaco-db App', () => {
  let page: PharmacoDBPage;

  beforeEach(() => {
    page = new PharmacoDBPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!');
  });
});
