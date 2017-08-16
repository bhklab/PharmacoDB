class CreateCellTissues < ActiveRecord::Migration[5.0]
	def change
		create_table :cell_tissues do |t|
			t.belongs_to :cell, index: true
			t.belongs_to :tissue, index: true

			t.timestamps
		end
	end
end
