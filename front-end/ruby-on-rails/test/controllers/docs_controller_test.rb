require 'test_helper'

class DocsControllerTest < ActionDispatch::IntegrationTest
  test "should get api" do
    get docs_api_url
    assert_response :success
  end

  test "should get cell_drug" do
    get docs_cell_drug_url
    assert_response :success
  end

  test "should get cells" do
    get docs_cells_url
    assert_response :success
  end

  test "should get tissues" do
    get docs_tissues_url
    assert_response :success
  end

  test "should get drugs" do
    get docs_drugs_url
    assert_response :success
  end

  test "should get datasets" do
    get docs_datasets_url
    assert_response :success
  end

  test "should get datasets_intersection" do
    get docs_datasets_intersection_url
    assert_response :success
  end

  test "should get search" do
    get docs_search_url
    assert_response :success
  end

end
