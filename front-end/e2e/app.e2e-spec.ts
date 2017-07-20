import { FrontEndPage } from './app.po';

describe('front-end App', () => {
  let page: FrontEndPage;

  beforeEach(() => {
    page = new FrontEndPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!!');
  });
});
