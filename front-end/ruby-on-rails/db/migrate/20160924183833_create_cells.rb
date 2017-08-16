class CreateCells < ActiveRecord::Migration[5.0]
	def change
		create_table :cells, id: false do |t|
			t.integer :cell_id
			t.string :accession_id
			t.text :cell_name

			t.timestamps
		end
		execute "ALTER TABLE cells ADD PRIMARY KEY (cell_id);"
	end
end
