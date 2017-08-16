class CreateSourceCellNames < ActiveRecord::Migration[5.0]
	def change
		create_table :source_cell_names do |t|
			t.belongs_to :cell, index: true
			t.belongs_to :source, index: true
			t.text :cell_name

			t.timestamps
		end
	end
end
