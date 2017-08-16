class CreateDatasets < ActiveRecord::Migration[5.0]
	def change
		create_table :datasets, id: false do |t|
			t.integer :dataset_id
			t.text :dataset_name

			t.timestamps
		end
		execute "ALTER TABLE datasets ADD PRIMARY KEY (dataset_id);"
	end
end
