class CreateCellosaurus < ActiveRecord::Migration[5.0]
  def change
    create_table :cellosaurus do |t|
      t.text :identifier
      t.text :accession_id
      t.text :synonyms
      t.text :creferences
      t.text :ireferences
      t.text :webpages
      t.text :comments
      t.text :strdata
      t.text :diseases
      t.text :sorigin
      t.text :hierarchy
      t.text :ofsi
      t.text :sgofcell
      t.text :category

      t.timestamps
    end
  end
end
