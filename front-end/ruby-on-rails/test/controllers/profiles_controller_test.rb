require 'test_helper'

class ProfilesControllerTest < ActionDispatch::IntegrationTest
  test "should get cell_lines" do
    get profiles_cell_lines_url
    assert_response :success
  end

  test "should get tissues" do
    get profiles_tissues_url
    assert_response :success
  end

  test "should get drugs" do
    get profiles_drugs_url
    assert_response :success
  end

  test "should get datasets" do
    get profiles_datasets_url
    assert_response :success
  end

end
